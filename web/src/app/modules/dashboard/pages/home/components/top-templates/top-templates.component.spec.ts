import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TopTemplateComponent } from './top-templates.component';

describe('TopTemplatesComponent', () => {
  let component: TopTemplateComponent;
  let fixture: ComponentFixture<TopTemplateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TopTemplateComponent]
    })
      .compileComponents();

    fixture = TestBed.createComponent(TopTemplateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
