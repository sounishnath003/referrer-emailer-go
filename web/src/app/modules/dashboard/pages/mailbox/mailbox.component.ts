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

  // Pagination & Filtering
  currentPage: number = 1;
  pageSize: number = 10;
  totalCount: number = 0;
  startDate: string = '';
  endDate: string = '';

  constructor(
    private emailingService: EmailingService,
    private profileService: ProfileService
  ) {
    this.searchSubject.pipe(
      debounceTime(300),
      distinctUntilChanged(),
      switchMap(query => {
        this.isLoading = true;
        this.currentPage = 1; // Reset to page 1 on search
        return this.emailingService.pollReferralMailbox$(
          this.profileService.ownerEmailAddress,
          query,
          this.currentPage,
          this.pageSize,
          this.startDate,
          this.endDate
        );
      })
    ).subscribe({
      next: (response) => {
        this.sentEmails = response.data || [];
        this.totalCount = response.meta.total;
        this.isLoading = false;
      },
      error: () => {
        this.isLoading = false;
        this.sentEmails = [];
        this.totalCount = 0;
      }
    });
  }

  ngOnInit() {
    this.loadEmails();
  }

  loadEmails() {
    this.isLoading = true;
    this.emailingService.pollReferralMailbox$(
      this.profileService.ownerEmailAddress,
      this.searchQuery,
      this.currentPage,
      this.pageSize,
      this.startDate,
      this.endDate
    ).subscribe({
      next: (response) => {
        this.sentEmails = response.data || [];
        this.totalCount = response.meta.total;
        this.isLoading = false;
      },
      error: () => {
        this.isLoading = false;
        this.sentEmails = [];
      }
    });
  }

  onSearch(query: string) {
    this.searchSubject.next(query);
  }

  onDateChange() {
    this.currentPage = 1;
    this.loadEmails();
  }

  nextPage() {
    if (this.currentPage * this.pageSize < this.totalCount) {
      this.currentPage++;
      this.loadEmails();
    }
  }

  prevPage() {
    if (this.currentPage > 1) {
      this.currentPage--;
      this.loadEmails();
    }
  }

  get totalPages(): number {
    return Math.ceil(this.totalCount / this.pageSize);
  }

  getCompanyName(subject: string): string {
    // TODO: Decide to implement or SKIP???
    // SKIP would be better option!
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

  stripHtml(html: string): string {
    const div = document.createElement('div');
    div.innerHTML = html;
    return div.textContent || div.innerText || '';
  }
}