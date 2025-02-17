import { Routes } from '@angular/router';
import { AuthComponent } from './auth.component';
import { LoginComponent } from './pages/login/login.component';
import { SignupComponent } from './pages/signup/signup.component';

export const routes: Routes = [
    {
        path: "",
        component: AuthComponent,
        children: [
            {
                path: "login",
                pathMatch: "full",
                component: LoginComponent
            },
            {
                path: "signup",
                pathMatch: "full",
                component: SignupComponent
            }
        ]
    }

];
