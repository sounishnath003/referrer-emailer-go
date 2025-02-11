import { AsyncPipe, NgIf } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { EmailingService } from '../../services/emailing.service';
import { Editor, NgxEditorModule } from 'ngx-editor';

@Component({
  selector: 'app-email-drafter',
  templateUrl: './email-drafter.component.html',
  styleUrls: ['./email-drafter.component.css'],
  imports: [FormsModule, ReactiveFormsModule, NgIf, AsyncPipe, NgxEditorModule],
  providers: [EmailingService]
})
export class EmailDrafterComponent implements OnInit, OnDestroy {
  editorBox!: Editor;
  html: string = "";

  toEmailIds: string[] = [];

  emailSenderForm: FormGroup = new FormGroup({
    to: new FormControl(null, [Validators.required, Validators.email]),
    from: new FormControl('flock.sinasini@gmai.com', [Validators.required, Validators.email]),
    subject: new FormControl(null, [Validators.required, Validators.maxLength(40)]),
    body: new FormControl(null, [Validators.required, Validators.minLength(30), Validators.maxLength(2000)]),
  });

  constructor(private emailingService: EmailingService) { }

  ngOnInit(): void {
    this.editorBox = new Editor();
  }
  ngOnDestroy(): void {
    this.editorBox.destroy();
    this.toEmailIds = [];
    this.emailSenderForm.reset();
  }

  // Function to handle email addition
  addEmail(): void {
    const emailControl = this.emailSenderForm.get('to');
    if (emailControl?.valid) {
      const email = emailControl.value.trim();
      if (!this.toEmailIds.includes(email)) {
        this.toEmailIds.push(email);
        emailControl.reset(); // Clear input after adding
      } else {
        window.alert('email exists')
      }
    }
  }

  // Remove email from the list
  removeEmail(email: string): void {
    this.toEmailIds = this.toEmailIds.filter((e) => e !== email);
  }

  onEmailSend() {
    console.log({ emailbox: this.html });

    const emailFormValue = this.emailSenderForm.value;
    // Check validations    
    if (this.toEmailIds.length == 0 || emailFormValue["body"].length < 10) {
      console.log(this.emailSenderForm.value);
      window.alert(JSON.stringify('Form is invalid'));
      return;
    }


    emailFormValue['to'] = [...this.toEmailIds, emailFormValue['from']];
    emailFormValue['body'] = emailFormValue['body'];

    // Call Api
    this.emailingService.sendEmail$(emailFormValue["from"], emailFormValue["to"], emailFormValue["subject"], emailFormValue["body"]).subscribe((resp) => {
      window.alert('Email has been sent');
      // Clear off
      this.emailSenderForm.reset();
      this.toEmailIds = [];
      console.log(resp);
    })

  }


  subscribeToFormUpdate$() {
    // return this.emailSenderForm.get('body')?.valueChanges;
    return this.emailSenderForm.valueChanges;
  }
}
