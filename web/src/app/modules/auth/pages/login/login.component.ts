import { CommonModule, NgStyle, TitleCasePipe } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { catchError, of } from 'rxjs';

@Component({
  selector: 'app-login',
  imports: [CommonModule, FormsModule, ReactiveFormsModule, RouterLink],
  providers: [AuthService],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})
export class LoginComponent implements OnInit {
  constructor(private readonly router: Router, private readonly authService: AuthService) { }
  loginForm: FormGroup = new FormGroup({
    email: new FormControl(null, [Validators.required, Validators.email]),
    password: new FormControl(null, [Validators.required, Validators.minLength(6)]),
  })
  formErrors: any = {};
  errorMessage: string | null = null;

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
    this.errorMessage = null;
    const errors: any = {};
    for (const controlName in this.loginForm.controls) {
      if (this.loginForm.controls[controlName].errors) {
        errors[controlName] = this.loginForm.controls[controlName].errors;
      }
    }
    return errors;
  }

  onLoginSubmit() {
    const { email, password } = this.loginForm.value;
    this.authService.login$(email, password).pipe(
      catchError(err => {
        this.errorMessage = err.error.error || `Not able to process requests. Try later`;
        return of(null);
      })
    ).subscribe({
      next: (data) => {
        if (data === null) return;
        const { accessToken } = data as { accessToken: string };
        // set the access token in local storage
        window.localStorage.setItem('REFERRER_ACCESS_TOKEN', accessToken);
        // move to dashboard page
        this.router.navigate(['dashboard'], { preserveFragment: false, });
      }
    })
  }

}
