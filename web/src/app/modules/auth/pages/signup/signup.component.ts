import { NgIf } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';

@Component({
  selector: 'app-signup',
  imports: [FormsModule, ReactiveFormsModule, RouterLink, NgIf],
  templateUrl: './signup.component.html',
  styleUrl: './signup.component.css'
})
export class SignupComponent implements OnInit {
  signupForm: FormGroup;
  formErrors: any = {};

  constructor(private readonly router: Router) {
    this.signupForm = new FormGroup({
      email: new FormControl(null, [Validators.required, Validators.email]),
      password: new FormControl(null, [Validators.required, Validators.minLength(6)]),
    })
  }

  ngOnInit(): void {
    this, this.signupForm.valueChanges.subscribe(() => this.onFormValueChange());
  }


  private onFormValueChange() {
    if (this.signupForm.invalid) {
      this.formErrors = this.getFormValidationErrors();
    } else {
      this.formErrors = {};
    }
  }

  private getFormValidationErrors() {
    const errors: any = {};
    for (const controlName in this.signupForm.controls) {
      if (this.signupForm.controls[controlName].errors) {
        errors[controlName] = this.signupForm.controls[controlName].errors;
      }
    }
    return errors;
  }


  onSignupSubmit() {
    this.router.navigate(['dashboard', 'profile'], { preserveFragment: false, skipLocationChange: false })
  }
}
