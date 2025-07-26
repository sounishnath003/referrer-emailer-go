import { Component, OnDestroy, OnInit } from '@angular/core';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { EmailingService } from '../../services/emailing.service';
import { Editor, NgxEditorModule } from 'ngx-editor';
import { EmailAutocompleteComponent } from "./components/email-autocomplete/email-autocomplete.component";
import { BehaviorSubject, catchError, debounceTime, distinctUntilChanged, filter, of, Subject, switchMap, takeUntil } from 'rxjs';
import { AsyncPipe } from '@angular/common';
import { SubheroComponent } from "../shared/subhero/subhero.component";
import { ApiProfileInformation, ProfileService } from '../../services/profile.service';

@Component({
  selector: 'app-email-drafter',
  templateUrl: './email-drafter.component.html',
  styleUrls: ['./email-drafter.component.css'],
  imports: [FormsModule, ReactiveFormsModule, NgxEditorModule, EmailAutocompleteComponent, AsyncPipe, SubheroComponent],
  providers: [EmailingService]
})
export class EmailDrafterComponent implements OnInit, OnDestroy {
  editorBox!: Editor;
  html: string = "";
  toEmailIds: string[] = [];
  filteredSuggestions: string[] = [];
  private destroy$ = new Subject<void>();
  processing$: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  successMessage: string | null = null;
  errorMessage: string | null = null;

  emailSenderForm: FormGroup = new FormGroup({
    to: new FormControl(null, [Validators.required, Validators.email]),
    from: new FormControl('sounish.nath17@gmail.com', [Validators.required, Validators.email]),
    subject: new FormControl(null, [Validators.required, Validators.maxLength(40)]),
    body: new FormControl(null, [Validators.required, Validators.minLength(30), Validators.maxLength(2000)]),
  });

  constructor(private readonly emailingService: EmailingService, private readonly profileService: ProfileService) { }

  ngOnInit(): void {
    this.editorBox = new Editor();

    this.profileService.getProfileInformation$(`sounish.nath17@gmail.com`)
      .pipe(
        catchError(err => {
          this.errorMessage = err.error.error || `Unable to fetch profile informations`;
          return of(null);
        })
      )
      .subscribe(data => {
        if (data == null) return;
        this.emailSenderForm.patchValue({
          subject: `Interested for [Job Profile] Roles Interview Opportunity at [Company Name]`,
          body: this.KickStartPrompt(data),
        }, { emitEvent: true })
      });

    // Update the Email form whenever the value changes
    this.emailSenderForm.get('to')?.valueChanges.pipe(
      debounceTime(300),
      distinctUntilChanged(),
      filter((value: string) => typeof value === 'string' && value.startsWith('@')),
      switchMap((value: string) => {
        const query = value.split('@').pop() || '';
        return this.profileService.searchPeople$(query);
      }),
      takeUntil(this.destroy$)
    ).subscribe({
      next: res => {
        this.filteredSuggestions = [...res.users];
      },
      error: err => {
        console.error(err);
        this.filteredSuggestions = [];
      }
    });


  }
  ngOnDestroy(): void {
    this.editorBox.destroy();
    this.toEmailIds = [];
    this.emailSenderForm.reset();

    this.destroy$.next();
    this.destroy$.complete();
  }

  onSuggestionSelected(suggestion: string | undefined): void {
    if (!suggestion) return;

    const inputControl = this.emailSenderForm.get('to');
    const currentValue = inputControl?.value || '';

    // Replace everything after the '@' with the selected suggestion
    const newValue = currentValue.replace(/@\w*$/, suggestion);
    inputControl?.setValue(newValue.trim());

    this.addEmail(); // move to chip
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


  private KickStartPrompt(user: ApiProfileInformation) {
    return `
  <p>Hello [User],</p>

  <p>Trust all is fine.</p>

  <p>It has always been my dream to work as a <b>[Job Role]</b> at <b>[Company Name]</b>, and I was thrilled to see a job opening that perfectly matches my skills and expertise.</p>

  <p>I kindly request that you review my resume and consider me for an interview if my qualifications align with your needs. Looking forward to hearing from you.</p>

  <p>For your convenience, I have attached my resume and the job post link below:</p>

  <ul>
    <li><b>Resume:</b> Attached in the email.</li>
    <li><b>Job Post:</b>
      <ul class="job-list">
        <li>Job Link #1</li>
        <li>Job Link #2</li>
        <li>Job Link #3</li>
      </ul>
    </li>
  </ul>

  <p>Thank you,</p>

  <div class="signature">
    <b>${user.firstName} ${user.lastName}</b>
    <p>${user.about}</p>
    <hr>
    <p>Email: ${user.email}</p>
  </div>
    `;
  }
}
