import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { ProfileInformation, ProfileService } from '../../services/profile.service';
import { RouterLink } from '@angular/router';
import { catchError, of } from 'rxjs';
import { SubheroComponent } from '../shared/subhero/subhero.component';
import { NgIf } from '@angular/common';

@Component({
  selector: 'app-profile',
  imports: [FormsModule, ReactiveFormsModule, RouterLink, SubheroComponent, NgIf],
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
      email: [{ value: "", disabled: true }],
      country: ['', Validators.required],
      currentCompany: [''],
      currentRole: [''],
      notifications: this.fb.group({
        receiveEmails: [true],
        offers: [false],
        pushNotifications: ['everything']
      })
    });
  }

  ngOnInit(): void {
    // get email from profileService
    const userEmail = this.profileService.ownerEmailAddress;
    this.profileService.getProfileInformation$(userEmail).pipe(
      catchError((err) => {
        this.errorMessage = err.error?.error || JSON.stringify(err.error);
        return of(null);
      })
    ).subscribe(
      data => {
        if (data === null) { return; }
        this.profileForm.patchValue({ ...this.profileForm.value, ...data }, { emitEvent: true, })
      }
    );
    this.profileForm.valueChanges.subscribe(() => this.onFormValueChange());
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
    const formValue: ProfileInformation = { ...this.profileForm.value, email: this.profileService.ownerEmailAddress } as ProfileInformation;

    this.profileService.updateProfileInformation$(formValue).pipe(
      catchError(err => {
        this.errorMessage = JSON.stringify(err.error.error);
        return of(null);
      })
    ).subscribe((data) => {
      if (data === null) {
        return;
      }
      this.errorMessage = null;
      this.successMessage = `Profile information has been updated.`;
      // Change the location to home page. in 2.5 seconds...
      setInterval(() => {
        window.location.replace("auth/login");
      }, 2500);
    });
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