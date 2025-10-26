import { Component, OnInit } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { SidebarViewComponent } from '../../shared/components/sidebar-view/sidebar-view.component';
import { MenuComponent } from "./pages/shared/menu/menu.component";
import { EmailingService, ReferralMailbox } from './services/emailing.service';
import { BehaviorSubject, catchError, of } from 'rxjs';
import { AsyncPipe } from '@angular/common';
import { ProfileService } from './services/profile.service';

@Component({
  selector: 'app-dashboard',
  imports: [RouterOutlet, SidebarViewComponent, MenuComponent, MenuComponent, AsyncPipe],
  providers: [EmailingService],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.css'
})
export class DashboardComponent implements OnInit {
  sentReferrals: BehaviorSubject<ReferralMailbox[]> = new BehaviorSubject<ReferralMailbox[]>([]);
  showSidebar: boolean = true;

  constructor(private readonly emailingService: EmailingService, private readonly profileService: ProfileService) { }

  ngOnInit(): void {
    this.pollReferralMailbox();
  }

  pollReferralMailbox(company?: string) {
    this.emailingService.pollReferralMailbox$(this.profileService.ownerEmailAddress, company).pipe(
      catchError(err => {
        console.error(err);
        this.sentReferrals.next([]);
        return of([]);
      })
    ).subscribe(
      data => {
        this.sentReferrals.next(data || []);
      }
    )
  }

  toggleSidebar() {
    this.showSidebar = !this.showSidebar;
  }
}
