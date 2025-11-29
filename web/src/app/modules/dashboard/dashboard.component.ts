import { Component, OnDestroy, OnInit } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { SidebarViewComponent } from '../../shared/components/sidebar-view/sidebar-view.component';
import { EmailingService } from './services/emailing.service';
import { ProfileService } from './services/profile.service';
import { Subject } from 'rxjs';

@Component({
  selector: 'app-dashboard',
  imports: [RouterOutlet, SidebarViewComponent],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.css'
})
export class DashboardComponent implements OnInit, OnDestroy {
  destroy$ = new Subject<void>();

  constructor(
    private readonly emailingService: EmailingService,
    private readonly profileService: ProfileService
  ) { }

  ngOnInit(): void {
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}
