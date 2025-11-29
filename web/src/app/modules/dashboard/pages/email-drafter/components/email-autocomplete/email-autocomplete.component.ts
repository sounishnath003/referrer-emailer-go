import { NgIf } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-email-autocomplete',
  imports: [NgIf],
  templateUrl: './email-autocomplete.component.html',
  styleUrl: './email-autocomplete.component.css'
})
export class EmailAutocompleteComponent {
  @Input() suggestions: { email: string, currentCompany: string }[] = [];
  @Output() suggestionSelected = new EventEmitter<{ email: string, currentCompany: string }>();

  selectSuggestion(suggestion: { email: string, currentCompany: string }) {
    this.suggestionSelected.emit(suggestion);
  }
}
