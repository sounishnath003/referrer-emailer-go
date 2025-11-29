import { Component, OnDestroy, OnInit } from '@angular/core';
import { SubheroComponent } from "../shared/subhero/subhero.component";
import { Editor, NgxEditorModule } from 'ngx-editor';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { AsyncPipe, CommonModule, JsonPipe, NgIf } from '@angular/common';
import { EmailingService } from '../../services/emailing.service';
import { BehaviorSubject, catchError, debounceTime, distinctUntilChanged, filter, of, Subject, switchMap, takeUntil } from 'rxjs';
import { MarkdownModule, MarkdownService } from 'ngx-markdown';
import { EmailAutocompleteComponent } from '../email-drafter/components/email-autocomplete/email-autocomplete.component';
import { ProfileService } from '../../services/profile.service';

@Component({
  selector: 'app-draft-with-ai',
  imports: [SubheroComponent, NgxEditorModule, FormsModule, ReactiveFormsModule, NgIf, MarkdownModule, AsyncPipe, EmailAutocompleteComponent, CommonModule],
  providers: [EmailingService, MarkdownService],
  templateUrl: './draft-with-ai.component.html',
  styleUrls: ['./draft-with-ai.component.css']
})
export class DraftWithAiComponent implements OnInit, OnDestroy {
  editorBox!: Editor;
  html: string = "";
  formErrors: any = {};
  template: string | null = null;

  filteredSuggestions: { email: string, currentCompany: string }[] = [];
  toEmailIds: string[] = [];
  private destroy$ = new Subject<void>();

  apiErrorMsg: string | null = null;
  loading: boolean = false;
  tailoredResumeId: string | null = null;

  emailReferralForm: FormGroup = new FormGroup({
    from: new FormControl(null, [Validators.required, Validators.email]),
    to: new FormControl(null, [Validators.email]),
    jobUrls: new FormControl(null, [Validators.required, Validators.pattern(/(https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z]{2,}(\.[a-zA-Z]{2,})(\.[a-zA-Z]{2,})?\/[a-zA-Z0-9]{2,}|((https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z]{2,}(\.[a-zA-Z]{2,})(\.[a-zA-Z]{2,})?)|(https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z0-9]{2,}\.[a-zA-Z0-9]{2,}\.[a-zA-Z0-9]{2,}(\.[a-zA-Z0-9]{2,})?/g)]),
    jobDescription: new FormControl(null, [Validators.maxLength(2000)]),
    companyName: new FormControl('', [Validators.required]),
    templateType: new FormControl({ value: this.template, disabled: true, }, [Validators.required]),
    subject: new FormControl('', []),
    body: new FormControl('', []),
  })

  get isFormValid(): boolean {
    // Form is valid if the base form is valid (ignoring the empty 'to' input if chips exist)
    // AND we have at least one recipient (either in chips or in the input)
    const hasRecipient = this.toEmailIds.length > 0 || !!this.emailReferralForm.get('to')?.value;
    return this.emailReferralForm.valid && hasRecipient;
  }

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
      this.tailoredResumeId = params.get('tailoredResumeId');
      const companyName = params.get('companyName');
      const to = params.get('to');

      this.emailReferralForm.patchValue({
        templateType: t,
        from: this.profileService.ownerEmailAddress,
        companyName: companyName,
        to: to
      }, { emitEvent: true });
      
      // Auto-add if 'to' present
      if (to) {
         this.addEmail();
      }
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
        this.filteredSuggestions = res.map(r => ({ email: r.email, currentCompany: r.currentCompany }));
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
    // Destroy the event cycl.
    this.destroy$.next();
    this.destroy$.complete();
  }

  onSuggestionSelected(suggestion: { email: string, currentCompany: string } | undefined): void {
    if (!suggestion) return;

    const inputControl = this.emailReferralForm.get('to');
    const currentValue = inputControl?.value || '';

    // Replace everything after the '@' with the selected suggestion
    const newValue = currentValue.replace(/@\w*$/, suggestion.email);
    inputControl?.setValue(newValue.trim());
    
    // Update company if not set or empty
    const currentCompany = this.emailReferralForm.get('companyName')?.value;
    if (!currentCompany && suggestion.currentCompany) {
        this.emailReferralForm.patchValue({ companyName: suggestion.currentCompany });
    }

    this.addEmail();
    this.filteredSuggestions = [];
  }

  // Function to handle email addition
  addEmail(): void {
    const emailControl = this.emailReferralForm.get('to');
    if (emailControl?.valid) {
      const email = emailControl.value.trim();
      if (!this.toEmailIds.includes(email)) {
        this.toEmailIds.push(email);
        emailControl.reset(); // Clear input after adding
      } else {
        // window.alert('email exists')
      }
    }
  }

  // Remove email from the list
  removeEmail(email: string): void {
    this.toEmailIds = this.toEmailIds.filter((e) => e !== email);
  }

  generateAiEmail() {
    // Logic to generate AI email and update previewBody
    // For AI generation, we use the first email if multiple, or just generate generic.
    // The prompt takes `to`. If multiple, we can send the first one as context or handle it.
    // Let's use the first one if available in `toEmailIds`, else the input value.
    
    let toEmail = this.emailReferralForm.get('to')?.value;
    if (this.toEmailIds.length > 0) {
        toEmail = this.toEmailIds[0];
    }
    
    const { from, companyName, jobUrls, jobDescription, templateType } = { ...this.emailReferralForm.value, templateType: this.template };
    // extract all the urls
    const jobUrlsExtract = this.extractAllUrls(jobUrls) || [jobUrls];

    this.loading = true;
    this.apiErrorMsg = null;
    this.successMessage = null;
    // Pass the existing tailoredResumeId to the service call
    this.emailService.generateAiDraftColdEmail$(from, toEmail, companyName, jobDescription, templateType, jobUrlsExtract, this.tailoredResumeId || undefined).pipe(
      catchError(err => {
        this.apiErrorMsg = err.error.error || `Unable to process your request!`;
        return of(null);
      })
    ).subscribe(data => {
      if (data === null) {
        this.loading = false;
        return;
      }
      this.html = this.markdownService.parse(data.mailBody) as string;
      this.emailReferralForm.patchValue({
        subject: data.mailSubject,
        body: this.html
      }, { emitEvent: true });
      // If the backend returns a new ID, use it. Otherwise, keep the existing one.
      this.tailoredResumeId = data.tailoredResumeId || this.tailoredResumeId;
      this.loading = false;
      this.apiErrorMsg = null;
    })

  }

  isBulkSending: boolean = false;

  sendEmail() {
    // Logic to send email
    const { from, subject, body } = { ...this.emailReferralForm.value };
    
    // Combine all recipients
    const recipients = [...this.toEmailIds];
    
    // Include current input value if it looks like an email and isn't already added
    const inputTo = this.emailReferralForm.get('to')?.value;
    if (inputTo && typeof inputTo === 'string' && inputTo.includes('@') && !recipients.includes(inputTo)) {
        recipients.push(inputTo);
    }
    
    // Add sender (backend requires it in the 'to' list for logic consistency, though it filters it out for bulk targeting)
    recipients.push(from);
    
    if (recipients.length <= 1) { // Only sender
         window.alert("Please add at least one recipient.");
         return;
    }

    this.processing$.next(true);
    this.isBulkSending = false;

    this.emailService.sendEmail$(from, recipients, subject, body, this.tailoredResumeId || undefined).pipe(
      catchError(err => {
        this.apiErrorMsg = err.error.error || `Something went wrong`;
        this.processing$.next(false);
        return of(null);
      })
    ).subscribe((data: any) => {
      if (data === null) return;
      this.errorMessage = null;

      // Check if it's a bulk response (202 Accepted) or single (200 OK)
      if (data.jobId) {
         // Start Polling
         this.isBulkSending = true;
         this.successMessage = data.message;
         this.pollJobStatus(data.jobId);
      } else {
         this.successMessage = `Email has been sent.`;
         this.processing$.next(false);
         // Reset
         this.toEmailIds = [];
         this.emailReferralForm.reset();
      }
    })
  }

  jobProgress: number = 0;
  jobStatus: string = '';
  isJobComplete: boolean = false;

  pollJobStatus(jobId: string) {
    const interval = setInterval(() => {
      this.emailService.getBulkEmailJobStatus$(jobId).subscribe(job => {
        this.jobStatus = job.status;
        this.jobProgress = Math.round((job.sentCount / job.totalRecipients) * 100);
        
        if (job.status === 'COMPLETED' || job.status === 'FAILED') {
          clearInterval(interval);
          this.processing$.next(false);
          this.isJobComplete = true;
          this.successMessage = job.status === 'COMPLETED' ? 'All emails sent successfully!' : 'Some emails failed to send.';
          
          if (job.status === 'COMPLETED') {
            //  this.emailReferralForm.reset(); 
             setTimeout(() => {
                 this.isJobComplete = false;
                 this.isBulkSending = false;
                 this.jobProgress = 0;
             }, 5000);
          }
        }
      });
    }, 1000); // Poll every 1 second
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
