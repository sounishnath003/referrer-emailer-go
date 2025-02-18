import { Component, OnDestroy, OnInit } from '@angular/core';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { EmailingService } from '../../services/emailing.service';
import { Editor, NgxEditorModule } from 'ngx-editor';
import { EmailAutocompleteComponent } from "./components/email-autocomplete/email-autocomplete.component";
import { BehaviorSubject, catchError, of } from 'rxjs';
import { AsyncPipe } from '@angular/common';

@Component({
  selector: 'app-email-drafter',
  templateUrl: './email-drafter.component.html',
  styleUrls: ['./email-drafter.component.css'],
  imports: [FormsModule, ReactiveFormsModule, NgxEditorModule, EmailAutocompleteComponent, AsyncPipe],
  providers: [EmailingService]
})
export class EmailDrafterComponent implements OnInit, OnDestroy {
  editorBox!: Editor;
  html: string = "";
  toEmailIds: string[] = [];
  suggestions: string[] = ['example1@example.com', 'flock.sinasini@gmail.com', 'example3@example.com']; // Example suggestions
  filteredSuggestions: string[] = [];
  processing$: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  successMessage: string | null = null;
  errorMessage: string | null = null;

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

  onEmailInput(event: Event): void {
    const input = (event.target as HTMLInputElement).value;
    if (input.startsWith('@')) {
      const prefix = input.split('@').pop() || "";
      this.filteredSuggestions = this.suggestions.filter(s => s.startsWith(prefix));
    } else {
      this.filteredSuggestions = [];
    }
  }

  onSuggestionSelected(suggestion: string | undefined): void {
    if (suggestion === undefined) return;
    const currentTo = this.emailSenderForm.get('to')?.value;
    const newTo = currentTo.replace(/@\w*$/, `${suggestion}`);
    this.emailSenderForm.get('to')?.setValue(newTo);
    this.addEmail();
    this.filteredSuggestions = [];
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
    if (this.toEmailIds.length == 0 || emailFormValue["subject"].length < 3 || emailFormValue["body"].length < 10) {
      console.log(this.emailSenderForm.value);
      window.alert(JSON.stringify('Form is invalid'));
      return;
    }


    emailFormValue['to'] = [...this.toEmailIds, emailFormValue['from']];
    emailFormValue['body'] = emailFormValue['body'];

    // Call Api
    this.processing$.next(true);
    this.emailingService.sendEmail$(emailFormValue["from"], emailFormValue["to"], emailFormValue["subject"], emailFormValue["body"]).pipe(
      catchError(err => {
        this.errorMessage = err.error.error;
        return of(null);
      })
    ).subscribe((resp) => {
      if (resp === null) return;
      this.errorMessage = null;
      this.successMessage = `Email has been sent.`;
      this.processing$.next(false);

      // Clear off
      this.emailSenderForm.reset();
      this.toEmailIds = [];
      console.log(resp);
    })

  }


  subscribeToFormUpdate$() {
    // reset the error | success messages.
    this.successMessage = null;
    this.errorMessage = null;
    this.processing$.next(false);
    // return this.emailSenderForm.get('body')?.valueChanges;
    return this.emailSenderForm.valueChanges;
  }
}
