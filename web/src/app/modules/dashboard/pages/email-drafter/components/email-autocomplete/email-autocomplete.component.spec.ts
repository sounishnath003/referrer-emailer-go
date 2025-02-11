import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EmailAutocompleteComponent } from './email-autocomplete.component';

describe('EmailAutocompleteComponent', () => {
  let component: EmailAutocompleteComponent;
  let fixture: ComponentFixture<EmailAutocompleteComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [EmailAutocompleteComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(EmailAutocompleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
