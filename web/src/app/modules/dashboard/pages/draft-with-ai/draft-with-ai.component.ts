import { Component, OnDestroy, OnInit } from '@angular/core';
import { SubheroComponent } from "../shared/subhero/subhero.component";
import { Editor, NgxEditorModule } from 'ngx-editor';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NgIf } from '@angular/common';

@Component({
  selector: 'app-draft-with-ai',
  imports: [SubheroComponent, NgxEditorModule, FormsModule, ReactiveFormsModule, NgIf],
  templateUrl: './draft-with-ai.component.html',
  styleUrls: ['./draft-with-ai.component.css']
})
export class DraftWithAiComponent implements OnInit, OnDestroy {
  editorBox!: Editor;
  html: string = "";
  formErrors: any = {};
  template: string | null = null;

  emailReferralForm: FormGroup = new FormGroup({
    from: new FormControl('flock.sinasini@gmail.com', [Validators.required, Validators.email]),
    to: new FormControl(null, [Validators.required, Validators.email]),
    jobUrls: new FormControl(null, [Validators.required, Validators.pattern(/(https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z]{2,}(\.[a-zA-Z]{2,})(\.[a-zA-Z]{2,})?\/[a-zA-Z0-9]{2,}|((https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z]{2,}(\.[a-zA-Z]{2,})(\.[a-zA-Z]{2,})?)|(https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z0-9]{2,}\.[a-zA-Z0-9]{2,}\.[a-zA-Z0-9]{2,}(\.[a-zA-Z0-9]{2,})?/g)]),
    jobDescription: new FormControl(null, []),
    templateType: new FormControl({ value: this.template, disabled: true, }, [Validators.required]),
    subject: new FormControl(null, [Validators.required]),
    body: new FormControl(null, []),
  })

  constructor(private readonly route: ActivatedRoute, private readonly router: Router) { }

  ngOnInit(): void {
    this.editorBox = new Editor();
    this.route.queryParamMap.subscribe(params => {
      const t = params.get('template');
      if (t === null || t.length === 0) {
        this.router.navigate(['dashboard']);
        return;
      }
      this.template = t;
      this.emailReferralForm.patchValue({ templateType: t });
    });
    this.emailReferralForm.valueChanges.subscribe(() => this.onFormValueChange());
  }

  ngOnDestroy(): void {
    this.editorBox.destroy();
  }

  generateAiEmail() {
    // Logic to generate AI email and update previewBody
    this.sendEmail();
  }

  sendEmail() {
    // Logic to send email
    const data = this.emailReferralForm.value;
    console.log(data);
  }


  private onFormValueChange() {
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
}