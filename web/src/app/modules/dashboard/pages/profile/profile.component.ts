import { NgIf } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { ProfileInformation, ProfileService } from '../../services/profile.service';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-profile',
  imports: [FormsModule, ReactiveFormsModule, NgIf, RouterLink],
  providers: [ProfileService],
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {
  profileForm: FormGroup;
  formErrors: any = {};
  errorMessage: string | null = null; // property to hanle API errors messages
  successMessage: string | null = null; // property to hanle API success messages

  constructor(private fb: FormBuilder, private readonly profileService: ProfileService) {
    this.profileForm = this.fb.group({
      firstName: ["", [Validators.required, Validators.minLength(3)]],
      lastName: ["", [Validators.required, Validators.minLength(3)]],
      about: ["", [Validators.required, Validators.minLength(50)]],
      resume: ["", [Validators.required]],
      email: [{ value: "flock.sinasini@gmail.com", disabled: true }],
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
    const formValue: ProfileInformation = { ...this.profileForm.value, email: "flock.sinasini@gmail.com" } as ProfileInformation;

    this.profileService.updateProfileInformation$(formValue).subscribe((data) => {
      this.errorMessage = null;
      this.successMessage = `Profile information has been updated.`;
    }, (err) => {
      this.errorMessage = JSON.stringify(err.error?.error || `Failed to update information. Please try again later.`);
      this.successMessage = null;
    })
  }

  private onFormValueChange() {
    if (this.profileForm.invalid) {
      this.formErrors = this.getFormValidationErrors();
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