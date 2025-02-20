import { Component } from '@angular/core';
import { RouterLink, RouterOutlet } from '@angular/router';
import { SidebarViewComponent } from '../../shared/components/sidebar-view/sidebar-view.component';
import { MenuComponent } from "./pages/shared/menu/menu.component";
import { SubheroComponent } from "./pages/shared/subhero/subhero.component";

@Component({
  selector: 'app-dashboard',
  imports: [RouterOutlet, SidebarViewComponent, MenuComponent, MenuComponent],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.css'
})
export class DashboardComponent {
  sentReferrals: number[] = [1, 2, 3];
}
