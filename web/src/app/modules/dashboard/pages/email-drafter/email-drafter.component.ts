import { Component } from '@angular/core';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';

@Component({
  selector: 'app-email-drafter',
  imports: [FormsModule, ReactiveFormsModule],
  templateUrl: './email-drafter.component.html',
  styleUrl: './email-drafter.component.css'
})
export class EmailDrafterComponent {
  emailSenderForm: FormGroup = new FormGroup({
    from: new FormControl(null, [Validators.required, Validators.email]),
    to: new FormControl(null, [Validators.required, Validators.email,]),
    subject: new FormControl(null, [Validators.required, Validators.maxLength(40),]),
    body: new FormControl(null, [Validators.required, Validators.maxLength(2000)]),
  })
}
