import { Component, OnInit } from '@angular/core';
import { ApiProfileInformation, ProfileService } from '../../services/profile.service';
import { catchError, Observable, of, tap } from 'rxjs';
import { AsyncPipe, NgIf } from '@angular/common';
import { MarkdownModule, MarkdownService } from 'ngx-markdown';
import { SubheroComponent } from '../shared/subhero/subhero.component';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-resume-view',
  imports: [AsyncPipe, NgIf, MarkdownModule, SubheroComponent, FormsModule],
  providers: [ProfileService, MarkdownService],
  templateUrl: './resume-view.component.html',
  styleUrl: './resume-view.component.css'
})
export class ResumeViewComponent implements OnInit {
  profileInformation$: Observable<ApiProfileInformation> | null = null;
  errorMessage: string | null = null;
  editMode = false;
  editedContent: string = '';

  constructor(private readonly profileService: ProfileService, private readonly markdownService: MarkdownService) { }

  ngOnInit(): void {
    this.profileInformation$ = this.profileService.getProfileInformation$(this.profileService.ownerEmailAddress).pipe(
      tap(profile => {
        if (profile) {
          this.editedContent = profile.extractedContent;
        }
      }),
      catchError(err => {
        this.errorMessage = err.error?.error || `No matches found.`;
        return of(null);
      }),
    ) as Observable<ApiProfileInformation>;
  }

  parseMarkdown(content: string) {
    return this.markdownService.parse(content);
  }

  save() {
    this.profileService.patchProfileInformation$({ extractedContent: this.editedContent, email: this.profileService.ownerEmailAddress }).subscribe({
      next: () => {
        this.editMode = false;
        // Refresh data
        this.ngOnInit();
      },
      error: (err) => {
        this.errorMessage = err.error?.error || 'Failed to save.';
      }
    });
  }
}
