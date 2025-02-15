import { Component } from '@angular/core';
import { RouterLink, RouterOutlet } from '@angular/router';
import { SidebarViewComponent } from '../../shared/components/sidebar-view/sidebar-view.component';

@Component({
  selector: 'app-dashboard',
  imports: [RouterOutlet, RouterLink, SidebarViewComponent],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.css'
})
export class DashboardComponent {
  sentReferrals: number[] = [1,2,3,4];
}
