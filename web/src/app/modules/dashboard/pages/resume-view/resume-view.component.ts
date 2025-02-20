import { Component, OnInit } from '@angular/core';
import { ApiProfileInformation, ProfileService } from '../../services/profile.service';
import { catchError, Observable, of } from 'rxjs';
import { AsyncPipe, NgIf } from '@angular/common';
import { MarkdownModule, MarkdownService } from 'ngx-markdown';
import { SubheroComponent } from '../shared/subhero/subhero.component';

@Component({
  selector: 'app-resume-view',
  imports: [AsyncPipe, NgIf, MarkdownModule, SubheroComponent],
  providers: [ProfileService, MarkdownService],
  templateUrl: './resume-view.component.html',
  styleUrl: './resume-view.component.css'
})
export class ResumeViewComponent implements OnInit {
  profileInformation$: Observable<ApiProfileInformation> | null = null;
  errorMessage: string | null = null;

  constructor(private readonly profileService: ProfileService, private readonly markdownService: MarkdownService) { }

  ngOnInit(): void {
    this.profileInformation$ = this.profileService.getProfileInformation$('flock.sinasini@gmail.com').pipe(
      catchError(err => {
        this.errorMessage = err.error?.error || `No matches found.`;
        return of(null);
      }),
    ) as Observable<ApiProfileInformation>;
  }

  parseMarkdown(content: string) {
    return this.markdownService.parse(content);
  }
}
