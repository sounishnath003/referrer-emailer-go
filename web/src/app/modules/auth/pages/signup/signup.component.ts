import { NgIf, TitleCasePipe } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { catchError, of } from 'rxjs';

@Component({
  selector: 'app-signup',
  imports: [FormsModule, ReactiveFormsModule, RouterLink, NgIf, TitleCasePipe],
  providers: [AuthService],
  templateUrl: './signup.component.html',
  styleUrl: './signup.component.css'
})
export class SignupComponent implements OnInit {
  signupForm: FormGroup;
  formErrors: any = {};
  errorMessage: string | null = null;

  constructor(private readonly router: Router, private readonly authService: AuthService) {
    this.signupForm = new FormGroup({
      email: new FormControl(null, [Validators.required, Validators.email]),
      password: new FormControl(null, [Validators.required, Validators.minLength(6)]),
    })
  }

  ngOnInit(): void {
    this, this.signupForm.valueChanges.subscribe(() => this.onFormValueChange());
  }


  private onFormValueChange() {
    this.errorMessage = null;
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
    const { email, password } = this.signupForm.value as { email: string, password: string };
    this.authService.signup$(email, password).pipe(
      catchError(err => {
        this.errorMessage = err.error.error;
        return of(null);
      })
    ).subscribe({
      next: (data) => {
        if (data === null) return;
        this.errorMessage = null;
        this.router.navigate(['dashboard', 'profile'], { preserveFragment: false, skipLocationChange: false, queryParams: { email } })
      },
    })
  }
}
