import { Routes } from '@angular/router';
import { DashboardComponent } from './dashboard.component';
import { EmailDrafterComponent } from './pages/email-drafter/email-drafter.component';
import { HomeComponent } from './pages/home/home.component';

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
        ]
    },
];
