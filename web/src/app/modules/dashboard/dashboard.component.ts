import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { SidebarViewComponent } from '../../shared/components/sidebar-view/sidebar-view.component';

@Component({
  selector: 'app-dashboard',
  imports: [RouterOutlet, SidebarViewComponent],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.css'
})
export class DashboardComponent {

}
