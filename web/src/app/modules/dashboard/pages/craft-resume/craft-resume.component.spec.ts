import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CraftResumeComponent } from './craft-resume.component';

describe('CraftResumeComponent', () => {
  let component: CraftResumeComponent;
  let fixture: ComponentFixture<CraftResumeComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CraftResumeComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CraftResumeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
