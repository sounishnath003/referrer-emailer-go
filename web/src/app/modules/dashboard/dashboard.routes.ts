import { Routes } from '@angular/router';
import { DashboardComponent } from './dashboard.component';
import { EmailDrafterComponent } from './pages/email-drafter/email-drafter.component';
import { HomeComponent } from './pages/home/home.component';
import { ProfileComponent } from './pages/profile/profile.component';
import { ResumeViewComponent } from './pages/resume-view/resume-view.component';

export const routes: Routes = [
    {
        path: "",
        component: DashboardComponent,
        children: [
            {
                path: "",
                pathMatch: "full",
                component: HomeComponent
            },
            {
                path: "email-drafter",
                pathMatch: "full",
                component: EmailDrafterComponent
            },
            {
                path: "profile",
                pathMatch: "full",
                component: ProfileComponent
            },
            {
                path: "resume",
                pathMatch: "full",
                component: ResumeViewComponent
            }
        ]
    },
];
