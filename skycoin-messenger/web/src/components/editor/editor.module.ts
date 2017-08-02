import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { EditorComponent } from './editor.component'

@NgModule({
  imports: [
    CommonModule
  ],
  exports: [
    EditorComponent,
  ],
  declarations: [EditorComponent]
})
export class EditorModule { }
