import { Routes } from '@angular/router';
import { DashboardComponent } from './dashboard.component';
import { EmailDrafterComponent } from './pages/email-drafter/email-drafter.component';
import { HomeComponent } from './pages/home/home.component';
import { ProfileComponent } from './pages/profile/profile.component';
import { ResumeViewComponent } from './pages/resume-view/resume-view.component';
import { SentReferralsComponent } from './pages/sent-referrals/sent-referrals.component';
import { DraftWithAiComponent } from './pages/draft-with-ai/draft-with-ai.component';

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
                path: "profile/update",
                pathMatch: "full",
                component: ProfileComponent
            },
            {
                path: "profile/resume",
                pathMatch: "full",
                component: ResumeViewComponent
            },
            {
                path: 'sent-referrals/:uuid',
                pathMatch: 'full',
                component: SentReferralsComponent,
            },
            {
                path: "draft-with-ai",
                pathMatch: "full",
                component: DraftWithAiComponent
            }
        ]
    },
];
