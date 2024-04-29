import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CspSpaComponent } from './csp-spa.component';

describe('CspSpaComponent', () => {
  let component: CspSpaComponent;
  let fixture: ComponentFixture<CspSpaComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CspSpaComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(CspSpaComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
