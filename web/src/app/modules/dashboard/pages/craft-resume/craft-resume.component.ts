import { Component } from '@angular/core';
import { SubheroComponent } from '../shared/subhero/subhero.component';
import { NgIf } from '@angular/common';
import { MarkdownModule } from 'ngx-markdown';
import { ProfileService } from '../../services/profile.service';
import { FormsModule } from '@angular/forms';
import { NgClass } from '@angular/common';
import { Router } from '@angular/router';

@Component({
  selector: 'app-craft-resume',
  imports: [SubheroComponent, NgIf, NgClass, MarkdownModule, FormsModule],
  templateUrl: './craft-resume.component.html',
  styleUrl: './craft-resume.component.css'
})
export class CraftResumeComponent {
  jobDescription = '';
  resumeMarkdown: string | null = null;
  loading = false;
  error: string | null = null;
  companyName = '';
  jobRole = '';

  constructor(private profileService: ProfileService, private router: Router) { }

  tailorResume(userEmail: string) {
    this.loading = true;
    this.error = null;
    this.resumeMarkdown = null;
    this.profileService.tailorResumeWithJobDescription$(this.jobDescription, userEmail, this.companyName, this.jobRole)
      .subscribe({
        next: (res) => {
          this.loading = false;
          if (res.id) {
            this.router.navigate([`/dashboard/craft-resume/${res.id}`]);
          }
        },
        error: (err) => {
          this.error = err.error?.error || 'Failed to generate tailored resume.';
          this.loading = false;
        }
      });
  }
}
