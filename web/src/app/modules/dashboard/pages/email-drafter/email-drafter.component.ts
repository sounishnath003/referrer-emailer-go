import { AsyncPipe, NgIf } from '@angular/common';
import { Component } from '@angular/core';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { MarkdownModule, MarkdownService } from 'ngx-markdown';
import { Observable } from 'rxjs';

@Component({
  selector: 'app-email-drafter',
  templateUrl: './email-drafter.component.html',
  styleUrls: ['./email-drafter.component.css'],
  imports: [FormsModule, ReactiveFormsModule, NgIf, AsyncPipe],
  providers: [MarkdownService]
})
export class EmailDrafterComponent {
  toEmailIds: string[] = [];

  emailSenderForm: FormGroup = new FormGroup({
    to: new FormControl(null, [Validators.required, Validators.email]),
    from: new FormControl(null, [Validators.required, Validators.email]),
    subject: new FormControl(null, [Validators.required, Validators.maxLength(40)]),
    body: new FormControl(null, [Validators.required, Validators.maxLength(2000)]),
  });

  constructor(private readonly markdownService: MarkdownService) { }

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
    window.alert('email has been sent')
  }

  parseMarkdownContent(content: string): string {
    console.log(content);

    return this.markdownService.parse(content) as string;
  }

  subscribeToFormUpdate$() {
    // return this.emailSenderForm.get('body')?.valueChanges;
    return this.emailSenderForm.valueChanges;
  }
}
