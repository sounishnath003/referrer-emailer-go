import { CommonModule, NgStyle } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';

@Component({
  selector: 'app-login',
  imports: [CommonModule, FormsModule, ReactiveFormsModule, RouterLink],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})
export class LoginComponent implements OnInit {
  constructor(private readonly router: Router) { }
  loginForm: FormGroup = new FormGroup({
    email: new FormControl(null, [Validators.required, Validators.email]),
    password: new FormControl(null, [Validators.required, Validators.minLength(6)]),
  })
  formErrors: any = {};

  ngOnInit() {
    this.loginForm.valueChanges.subscribe(() => {
      this.onFormValueChange();
    });
  }

  onFormValueChange() {
    if (this.loginForm.invalid) {
      this.formErrors = this.getFormValidationErrors();
    } else {
      this.formErrors = {};
    }
  }

  getFormValidationErrors() {
    const errors: any = {};
    for (const controlName in this.loginForm.controls) {
      if (this.loginForm.controls[controlName].errors) {
        errors[controlName] = this.loginForm.controls[controlName].errors;
      }
    }
    return errors;
  }

  onLoginSubmit() {
    const value = this.loginForm.value;
    this.router.navigate(['dashboard'], { preserveFragment: false });
  }

}
