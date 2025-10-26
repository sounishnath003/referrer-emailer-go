import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { ReferralMailbox } from '../../../modules/dashboard/services/emailing.service';
import { CommonModule, DatePipe } from '@angular/common';
import { RouterLink } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { Subject } from 'rxjs';
import { debounceTime, distinctUntilChanged } from 'rxjs/operators';

@Component({
  selector: 'app-sidebar-view',
  imports: [RouterLink, CommonModule, FormsModule, DatePipe],
  templateUrl: './sidebar-view.component.html',
  styleUrl: './sidebar-view.component.css'
})
export class SidebarViewComponent implements OnInit, OnDestroy {
  @Input() sentReferrals!: ReferralMailbox[];
  @Output() searchTermChange = new EventEmitter<string>();

  private searchDebounce = new Subject<string>();
  searchTerm: string = '';

  constructor() { }

  ngOnInit(): void {
    this.searchDebounce.pipe(
      debounceTime(300), // 300ms debounce time
      distinctUntilChanged()
    ).subscribe(term => {
      this.searchTermChange.emit(term);
    });
  }

  ngOnDestroy(): void {
    this.searchDebounce.unsubscribe();
  }

  onSearchTermChange(term: string): void {
    this.searchDebounce.next(term);
  }

  getEmailPrefix(email: string): string {
    if (!email) {
      return '';
    }
    const atIndex = email.indexOf('@');
    return atIndex !== -1 ? email.substring(0, atIndex) : email;
  }

  getCompanyName(subject: string): string {
    if (!subject) {
      return 'N/A';
    }
    const patterns = [
      /Interested for (.*?)(?: -|$)/i,
      /Application for (.*?)(?: at|$)/i,
      /Referral for (.*?)(?: at|$)/i,
    ];
    for (const pattern of patterns) {
      const match = subject.match(pattern);
      if (match && match[1]) {
        return match[1].trim();
      }
    }
    return subject.split(/for|at|-/)[0].trim() || subject;
  }
}
