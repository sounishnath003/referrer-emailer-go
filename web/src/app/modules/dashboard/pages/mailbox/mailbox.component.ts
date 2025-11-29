import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { EmailingService, ReferralMailbox } from '../../services/emailing.service';
import { ProfileService } from '../../services/profile.service';
import { FormsModule } from '@angular/forms';
import { Subject, debounceTime, distinctUntilChanged, switchMap } from 'rxjs';

@Component({
  selector: 'app-mailbox',
  standalone: true,
  imports: [CommonModule, RouterModule, FormsModule],
  templateUrl: './mailbox.component.html',
  styleUrl: './mailbox.component.css'
})
export class MailboxComponent implements OnInit {
  sentEmails: ReferralMailbox[] = [];
  isLoading: boolean = false;
  searchQuery: string = '';
  searchSubject = new Subject<string>();
  
  constructor(
    private emailingService: EmailingService,
    private profileService: ProfileService
  ) {
    this.searchSubject.pipe(
      debounceTime(300),
      distinctUntilChanged(),
      switchMap(query => {
        this.isLoading = true;
        return this.emailingService.pollReferralMailbox$(this.profileService.ownerEmailAddress, query);
      })
    ).subscribe({
      next: (emails) => {
        this.sentEmails = emails;
        this.isLoading = false;
      },
      error: () => {
        this.isLoading = false;
        this.sentEmails = [];
      }
    });
  }

  ngOnInit() {
    this.loadEmails();
  }

  loadEmails() {
    this.isLoading = true;
    this.emailingService.pollReferralMailbox$(this.profileService.ownerEmailAddress, this.searchQuery).subscribe({
      next: (emails) => {
        this.sentEmails = emails;
        this.isLoading = false;
      },
      error: () => {
        this.isLoading = false;
      }
    });
  }

  onSearch(query: string) {
    this.searchSubject.next(query);
  }

  getCompanyName(subject: string): string {
    // Basic extraction if format is consistent, otherwise return subject
    // Example: "Interested for [Role] at [Company]"
    // Or just return the whole subject as "Subject" column
    return subject;
  }

  getEmailDisplay(to: string[]): string {
    if (!to || to.length === 0) return 'No Recipient';
    // Filter out owner email if possible, or just show first
    const owner = this.profileService.ownerEmailAddress;
    const recipients = to.filter(e => e.toLowerCase() !== owner.toLowerCase());
    if (recipients.length > 0) return recipients.join(', ');
    return to.join(', ');
  }
}