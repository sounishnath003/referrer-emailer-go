import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-sidebar-view',
  imports: [],
  templateUrl: './sidebar-view.component.html',
  styleUrl: './sidebar-view.component.css'
})
export class SidebarViewComponent {
  @Input() sentReferrals!: number[];
}
