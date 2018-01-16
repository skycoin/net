import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { DashboardComponent, SubStatusComponent, LoginComponent, UpdatePassComponent } from '../page';

const routes: Routes = [
  {
    path: '',
    component: DashboardComponent,
    pathMatch: 'full'
  },
  {
    path: 'node',
    component: SubStatusComponent
  },
  {
    path: 'login',
    component: LoginComponent,
    outlet: 'user'
  },
  {
    path: 'updatePass',
    component: UpdatePassComponent,
  },
  { path: '**', redirectTo: '' },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { useHash: true })],
  exports: [RouterModule],
})
export class AppRouteModule { }
