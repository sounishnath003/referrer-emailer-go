import { Component, OnDestroy, OnInit } from '@angular/core';
import { SubheroComponent } from "../shared/subhero/subhero.component";
import { Editor, NgxEditorModule } from 'ngx-editor';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { AsyncPipe, JsonPipe, NgIf } from '@angular/common';
import { EmailingService } from '../../services/emailing.service';
import { BehaviorSubject, catchError, debounceTime, distinctUntilChanged, filter, of, Subject, switchMap, takeUntil } from 'rxjs';
import { MarkdownModule, MarkdownService } from 'ngx-markdown';
import { EmailAutocompleteComponent } from '../email-drafter/components/email-autocomplete/email-autocomplete.component';
import { ProfileService } from '../../services/profile.service';

@Component({
  selector: 'app-draft-with-ai',
  imports: [SubheroComponent, NgxEditorModule, FormsModule, ReactiveFormsModule, NgIf, MarkdownModule, AsyncPipe, EmailAutocompleteComponent],
  providers: [EmailingService, MarkdownService],
  templateUrl: './draft-with-ai.component.html',
  styleUrls: ['./draft-with-ai.component.css']
})
export class DraftWithAiComponent implements OnInit, OnDestroy {
  editorBox!: Editor;
  html: string = "";
  formErrors: any = {};
  template: string | null = null;

  filteredSuggestions: { email: string, companyName: string }[] = [];
  private destroy$ = new Subject<void>();

  apiErrorMsg: string | null = null;
  loading: boolean = false;
  tailoredResumeId: string | null = null;

  emailReferralForm: FormGroup = new FormGroup({
    from: new FormControl(null, [Validators.required, Validators.email]),
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


  constructor(private readonly route: ActivatedRoute, private readonly router: Router, private readonly emailService: EmailingService, private readonly markdownService: MarkdownService, private readonly profileService: ProfileService) { }

  ngOnInit(): void {
    this.editorBox = new Editor();
    this.route.queryParamMap.subscribe(params => {
      const t = params.get('template');
      if (t === null || t.length === 0) {
        this.router.navigate(['dashboard']);
        return;
      }
      this.template = t;
      this.emailReferralForm.patchValue({ templateType: t, from: this.profileService.ownerEmailAddress }, { emitEvent: true });
    });
    this.emailReferralForm.valueChanges.subscribe(() => this.onFormValueChange());

    // Update the Email form whenever the value changes
    this.emailReferralForm.get('to')?.valueChanges.pipe(
      debounceTime(300),
      distinctUntilChanged(),
      filter((value: string) => typeof value === 'string' && value.startsWith('@')),
      switchMap((value: string) => {
        const query = value.trim().split('@').pop() || '';
        return this.profileService.searchPeople$(query);
      }),
      takeUntil(this.destroy$)
    ).subscribe({
      next: res => {
        this.filteredSuggestions = [...res];
      },
      error: err => {
        console.error(err);
        this.filteredSuggestions = [];
      }
    });
  }

  ngOnDestroy(): void {
    this.editorBox.destroy();
    // Destroy the event cycl.
    this.destroy$.next();
    this.destroy$.complete();
  }

  onSuggestionSelected(suggestion: { email: string, companyName: string } | undefined): void {
    if (!suggestion) return;

    // PatchValues
    this.emailReferralForm.patchValue({
      to: suggestion.email.trim(),
      companyName: suggestion.companyName.trim(),
    });
    this.filteredSuggestions = [];
  }

  generateAiEmail() {
    // Logic to generate AI email and update previewBody
    const { from, to, companyName, jobUrls, jobDescription, templateType } = { ...this.emailReferralForm.value, templateType: this.template };
    // extract all the urls
    const jobUrlsExtract = this.extractAllUrls(jobUrls) || [jobUrls];

    this.loading = true;
    this.apiErrorMsg = null;
    this.successMessage = null;
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
      this.tailoredResumeId = data.tailoredResumeId || null;
      this.loading = false;
      this.apiErrorMsg = null;
    })

  }

  sendEmail() {
    // Logic to send email
    const { from, to, subject, body } = { ...this.emailReferralForm.value };

    this.processing$.next(true);

    this.emailService.sendEmail$(from, [to, from], subject, body, this.tailoredResumeId || undefined).pipe(
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
      // this.emailReferralForm.reset();
    })
  }


  private onFormValueChange() {
    this.apiErrorMsg = null;
    this.successMessage = null;

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