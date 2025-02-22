import { Component, OnDestroy, OnInit } from '@angular/core';
import { SubheroComponent } from "../shared/subhero/subhero.component";
import { Editor, NgxEditorModule } from 'ngx-editor';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { AsyncPipe, JsonPipe, NgIf } from '@angular/common';
import { EmailingService } from '../../services/emailing.service';
import { BehaviorSubject, catchError, of } from 'rxjs';
import { MarkdownModule, MarkdownService } from 'ngx-markdown';

@Component({
  selector: 'app-draft-with-ai',
  imports: [SubheroComponent, NgxEditorModule, FormsModule, ReactiveFormsModule, NgIf, MarkdownModule, AsyncPipe],
  providers: [EmailingService, MarkdownService],
  templateUrl: './draft-with-ai.component.html',
  styleUrls: ['./draft-with-ai.component.css']
})
export class DraftWithAiComponent implements OnInit, OnDestroy {
  editorBox!: Editor;
  html: string = "";
  formErrors: any = {};
  template: string | null = null;

  apiErrorMsg: string | null = null;
  loading: boolean = false;

  emailReferralForm: FormGroup = new FormGroup({
    from: new FormControl('flock.sinasini@gmail.com', [Validators.required, Validators.email]),
    to: new FormControl(null, [Validators.required, Validators.email]),
    jobUrls: new FormControl(null, [Validators.required, Validators.pattern(/(https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z]{2,}(\.[a-zA-Z]{2,})(\.[a-zA-Z]{2,})?\/[a-zA-Z0-9]{2,}|((https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z]{2,}(\.[a-zA-Z]{2,})(\.[a-zA-Z]{2,})?)|(https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z0-9]{2,}\.[a-zA-Z0-9]{2,}\.[a-zA-Z0-9]{2,}(\.[a-zA-Z0-9]{2,})?/g)]),
    jobDescription: new FormControl(null, [Validators.maxLength(2000)]),
    companyName: new FormControl('', [Validators.required]),
    templateType: new FormControl({ value: this.template, disabled: true, }, [Validators.required]),
    subject: new FormControl('', []),
    body: new FormControl('', []),
  })

  processing$: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  successMessage: string | null = null;
  errorMessage: string | null = null;


  constructor(private readonly route: ActivatedRoute, private readonly router: Router, private readonly emailService: EmailingService, private readonly markdownService: MarkdownService) { }

  ngOnInit(): void {
    this.editorBox = new Editor();
    this.route.queryParamMap.subscribe(params => {
      const t = params.get('template');
      if (t === null || t.length === 0) {
        this.router.navigate(['dashboard']);
        return;
      }
      this.template = t;
      this.emailReferralForm.patchValue({ templateType: t }, { emitEvent: true });
    });
    this.emailReferralForm.valueChanges.subscribe(() => this.onFormValueChange());
  }

  ngOnDestroy(): void {
    this.editorBox.destroy();
  }

  generateAiEmail() {
    // Logic to generate AI email and update previewBody
    const { from, to, companyName, jobUrls, jobDescription, templateType } = { ...this.emailReferralForm.value, templateType: this.template };
    // extract all the urls
    const jobUrlsExtract = this.extractAllUrls(jobUrls) || [jobUrls];

    this.emailService.generateAiDraftColdEmail$(from, to, companyName, jobDescription, templateType, jobUrlsExtract).pipe(
      catchError(err => {
        this.apiErrorMsg = err.error.error || `Unable to process your request!`;
        return of(null);
      })
    ).subscribe(data => {
      if (data === null) return;
      this.html = this.markdownService.parse(data.mailBody) as string;

      this.emailReferralForm.patchValue({
        subject: data.mailSubject,
        body: this.html
      }, { emitEvent: true });

      this.apiErrorMsg = null;
    })

  }

  sendEmail() {
    // Logic to send email
    const { from, to, subject, body } = { ...this.emailReferralForm.value };

    this.processing$.next(true);

    this.emailService.sendEmail$(from, [to, from], subject, body).pipe(
      catchError(err => {
        this.apiErrorMsg = err.error.error || `Something went wrong`;
        return of(null);
      })
    ).subscribe(data => {
      if (data === null) return;
      this.errorMessage = null;
      this.successMessage = `Email has been sent.`;
      this.processing$.next(false);

      // Clear off
      this.emailReferralForm.reset();
    })
  }


  private onFormValueChange() {
    this.apiErrorMsg = null;
    if (this.emailReferralForm.invalid) {
      this.formErrors = this.getFormValidationErrors();
    } else {
      this.formErrors = {};
    }
  }

  private getFormValidationErrors() {
    const errors: any = {};
    for (const controlName in this.emailReferralForm.controls) {
      if (this.emailReferralForm.controls[controlName].errors) {
        errors[controlName] = this.emailReferralForm.controls[controlName].errors;
      }
    }
    return errors;
  }

  private extractAllUrls(content: string) {
    const urlRegex = /(https?:\/\/[^\s]+)/g;
    const urls = content.match(urlRegex);
    return urls;
  }
}