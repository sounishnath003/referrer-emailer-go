import { Component, Input } from '@angular/core';
import { ReferralMailbox } from '../../../modules/dashboard/services/emailing.service';
import { TitleCasePipe } from '@angular/common';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-sidebar-view',
  imports: [TitleCasePipe, RouterLink],
  templateUrl: './sidebar-view.component.html',
  styleUrl: './sidebar-view.component.css'
})
export class SidebarViewComponent {
  @Input() sentReferrals!: ReferralMailbox[];

  constructor() { }

  parseMarkdown(content: string) {
    return content.slice(0, 150);
  }
}
