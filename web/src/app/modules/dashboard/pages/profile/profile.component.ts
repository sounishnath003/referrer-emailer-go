import { Component } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';

@Component({
  selector: 'app-profile',
  imports: [FormsModule, ReactiveFormsModule],
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent {
  profileForm: FormGroup;

  constructor(private fb: FormBuilder) {
    this.profileForm = this.fb.group({
      firstName: ['', Validators.required],
      lastName: ['', Validators.required],
      about: ['', [Validators.required, Validators.minLength(10), Validators.maxLength(200)]],
      resume: [null, Validators.required],
      email: [{ value: 'flock.sinasini@gmail.com', disabled: true }],
      country: ['', Validators.required],
      notifications: this.fb.group({
        comments: [true],
        candidates: [false],
        offers: [false],
        pushNotifications: ['everything']
      })
    });
  }

  onFileChange(event: any) {
    const file = event.target.files[0];
    if (file) {
      this.profileForm.patchValue({
        resume: file
      });
    }
  }

  onSubmit() {
    // if (this.profileForm.valid) {
    const formData = new FormData();
    Object.keys(this.profileForm.controls).forEach(key => {
      if (key === 'resume') {
        formData.append(key, this.profileForm.get(key)?.value);
      } else {
        formData.append(key, this.profileForm.get(key)?.value);
      }
    });
    // Handle form submission, e.g., send formData to the server
    console.log('Form submitted', formData);
    // }
  }
}