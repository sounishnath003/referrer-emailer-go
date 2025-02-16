import { AsyncPipe, JsonPipe, NgIf } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';

@Component({
  selector: 'app-profile',
  imports: [FormsModule, ReactiveFormsModule, NgIf, JsonPipe],
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {
  profileForm: FormGroup;
  formErrors: any = {};

  constructor(private fb: FormBuilder) {
    this.profileForm = this.fb.group({
      firstName: ["", [Validators.required, Validators.minLength(3)]],
      lastName: ["", [Validators.required, Validators.minLength(3)]],
      about: ["", [Validators.required, Validators.minLength(50)]],
      resume: ["", [Validators.required]],
      email: [{ value: 'flock.sinasini@gmail.com', disabled: true }],
      country: ['', Validators.required],
      notifications: this.fb.group({
        receiveEmails: [true],
        offers: [false],
        pushNotifications: ['everything']
      })
    });
  }

  ngOnInit(): void {
    this.profileForm.valueChanges.subscribe(() => this.onFormValueChange())
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
    console.log({ isValid: this.profileForm.valid, profileForm: this.profileForm.value });
  }

  private onFormValueChange() {
    if (this.profileForm.invalid) {
      this.formErrors = this.getFormValidationErrors();
      console.log(this.formErrors);

    } else {
      this.formErrors = {};
    }
  }

  private getFormValidationErrors() {
    const errors: any = {};
    for (const controlName in this.profileForm.controls) {
      if (this.profileForm.controls[controlName].errors) {
        errors[controlName] = this.profileForm.controls[controlName].errors;
      }
    }
    return errors;
  }
}