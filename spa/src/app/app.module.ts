import { NgModule } from '@angular/core'; 
import { BrowserModule } 
	from '@angular/platform-browser'; 
import { FormsModule } from '@angular/forms'; 
import { AppRoutingModule } 
	from './app.routes'; 
import { AppComponent } 
	from './app.component'; 
import { CspSpaComponent } 
	from './csp-spa/csp-spa.component'; 

@NgModule({ 
	declarations: [ 
		AppComponent, 
		CspSpaComponent 
	], 
	imports: [ 
		BrowserModule, 
		AppRoutingModule, 
		FormsModule 
	], 
	providers: [], 
	bootstrap: [AppComponent] 
}) 
export class AppModule { }
