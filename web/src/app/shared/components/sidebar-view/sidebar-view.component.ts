import { Component, Input } from '@angular/core';
import { ReferralMailbox } from '../../../modules/dashboard/services/emailing.service';
import { TitleCasePipe } from '@angular/common';
import { MarkdownModule, MarkdownPipe, MarkdownService } from 'ngx-markdown';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-sidebar-view',
  imports: [TitleCasePipe, MarkdownModule, RouterLink],
  providers: [MarkdownPipe],
  templateUrl: './sidebar-view.component.html',
  styleUrl: './sidebar-view.component.css'
})
export class SidebarViewComponent {
  @Input() sentReferrals!: ReferralMailbox[];

  constructor(private readonly markdownService: MarkdownService) { }

  parseMarkdown(content: string) {
    const updatedContent = this.markdownService.parse(content) as string;
    return `${updatedContent.slice(0, 70)}...`;
  }
}
