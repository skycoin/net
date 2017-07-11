import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FlexLayoutModule } from '@angular/flex-layout';
import { AppComponent } from './app.component';
import {
  ImViewComponent,
  ImRecentBarComponent,
  ImRecentItemComponent,
  ImHeadComponent,
  ImHistoryViewComponent,
  ImHistoryMessageComponent
} from '../components';
import { FormsModule } from '@angular/forms';
import { SocketService } from '../providers';


@NgModule({
  declarations: [
    AppComponent,
    ImViewComponent,
    ImRecentBarComponent,
    ImRecentItemComponent,
    ImHeadComponent,
    ImHistoryViewComponent,
    ImHistoryMessageComponent
  ],
  imports: [
    FormsModule,
    BrowserModule,
    FlexLayoutModule,
  ],
  providers: [SocketService],
  bootstrap: [AppComponent]
})
export class AppModule { }
