import { Routes } from '@angular/router';
import { DashboardComponent } from './dashboard.component';
import { EmailDrafterComponent } from './pages/email-drafter/email-drafter.component';
import { HomeComponent } from './pages/home/home.component';
import { ProfileComponent } from './pages/profile/profile.component';
import { ResumeViewComponent } from './pages/resume-view/resume-view.component';
import { SentReferralsComponent } from './pages/sent-referrals/sent-referrals.component';
import { DraftWithAiComponent } from './pages/draft-with-ai/draft-with-ai.component';
import { ProfileAnalyticsComponent } from './pages/profile-analytics/profile-analytics.component';
import { CraftResumeComponent } from './pages/craft-resume/craft-resume.component';
import { TailoredResumeViewComponent } from './pages/craft-resume/tailored-resume-view.component';
import { TailoredResumeListComponent } from './pages/craft-resume/tailored-resume-list.component';
import { NetworkComponent } from './pages/network/network.component';

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
                path: "network",
                pathMatch: "full",
                component: NetworkComponent
            },
            {
                path: "email-drafter",
                pathMatch: "full",
                component: EmailDrafterComponent
            },
            {
                path: "craft-resume",
                pathMatch: "full",
                component: CraftResumeComponent
            },
            {
                path: "craft-resume/:id",
                pathMatch: "full",
                component: TailoredResumeViewComponent
            },
            {
                path: "craft-resume-list",
                pathMatch: "full",
                component: TailoredResumeListComponent
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
                path: "profile/analytics",
                pathMatch: "full",
                component: ProfileAnalyticsComponent
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
