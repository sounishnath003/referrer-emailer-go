import { Component, OnInit } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { SidebarViewComponent } from '../../shared/components/sidebar-view/sidebar-view.component';
import { MenuComponent } from "./pages/shared/menu/menu.component";
import { EmailingService, ReferralMailbox } from './services/emailing.service';
import { BehaviorSubject, catchError, of } from 'rxjs';
import { AsyncPipe, NgIf } from '@angular/common';

@Component({
  selector: 'app-dashboard',
  imports: [RouterOutlet, SidebarViewComponent, MenuComponent, MenuComponent, AsyncPipe, NgIf],
  providers: [EmailingService],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.css'
})
export class DashboardComponent implements OnInit {
  sentReferrals: BehaviorSubject<ReferralMailbox[]> = new BehaviorSubject<ReferralMailbox[]>([]);
  constructor(private readonly emailingService: EmailingService) { }


  ngOnInit(): void {
    this.pollReferralMailbox();
  }

  pollReferralMailbox() {
    this.emailingService.pollReferralMailbox$(`flock.sinasini@gmail.com`).pipe(
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
}
