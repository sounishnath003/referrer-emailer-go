import { Routes } from '@angular/router';
import { Error404Component } from './shared/components/error-404/error-404.component';

export const routes: Routes = [
    {
        path: "",
        pathMatch: "full",
        redirectTo: "auth/login"
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
        path: "**",
        component: Error404Component
    }
];
