import { Routes } from '@angular/router';
import { AuthComponent } from './auth.component';
import { LoginComponent } from './pages/login/login.component';

export const routes: Routes = [
    {
        path: "",
        component: AuthComponent,
        children: [
            {
                path: "login",
                pathMatch: "full",
                component: LoginComponent
            }
        ]
    }

];
