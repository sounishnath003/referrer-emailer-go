import { AsyncPipe, NgFor, NgIf } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { ProfileAnalytics, ProfileService } from '../../services/profile.service';
import { catchError, Observable, of } from 'rxjs';
import { SubheroComponent } from "../shared/subhero/subhero.component";

@Component({
  selector: 'app-profile-analytics',
  imports: [NgFor, AsyncPipe, NgIf, SubheroComponent],
  providers: [ProfileService],
  templateUrl: './profile-analytics.component.html',
  styleUrl: './profile-analytics.component.css'
})
export class ProfileAnalyticsComponent implements OnInit {
  profileAnalytics$: Observable<ProfileAnalytics | null> | null = null;
  apiError: string | null = null;

  constructor(private readonly profileService: ProfileService) { }

  ngOnInit(): void {
    this.profileAnalytics$ = this.profileService.getProfileAnalytics$(`flock.sinasini@gmail.com`).pipe(
      catchError(err => {
        this.apiError = err.error.error || `Something went wrong. No response from backend`;
        return of(null);
      })
    );
  }
}
