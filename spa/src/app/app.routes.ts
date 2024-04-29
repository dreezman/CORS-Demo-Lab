import { NgModule } from '@angular/core'; 
import { Routes, RouterModule } 
	from '@angular/router'; 
import { CspSpaComponent } 
	from './csp-spa/csp-spa.component'; 

const routes: Routes = [ 
	{ path: '', component: CspSpaComponent } 
]; 

@NgModule({ 
  imports: [RouterModule.forRoot(routes)], 
  exports: [RouterModule] 
}) 
export class AppRoutingModule { }
