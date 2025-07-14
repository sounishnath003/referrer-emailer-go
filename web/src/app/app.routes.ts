import { Routes } from '@angular/router';
import { Error404Component } from './shared/components/error-404/error-404.component';
import { LandingComponent } from './modules/landing/landing.component';
import { TailoredResumeViewComponent } from './modules/dashboard/pages/craft-resume/tailored-resume-view.component';

export const routes: Routes = [
    {
        path: "",
        pathMatch: "full",
        redirectTo: "demo"
    },
    {
        path: "demo",
        pathMatch: "full",
        component: LandingComponent
    },
    {
        path: "auth",
        loadChildren: async () => (await import('./modules/auth/auth.module')).AuthModule
    },
    {
        path: "dashboard",
        loadChildren: async () => (await import('./modules/dashboard/dashboard.module')).DashboardModule
    },
    {
        path: "craft-resume/:id",
        component: TailoredResumeViewComponent
    },
    {
        path: "**",
        component: Error404Component
    }
];
