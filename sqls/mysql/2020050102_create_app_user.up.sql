create table `app_user` (
 `id` int auto_increment
,`version` int not null
,`created_at` datetime not null default current_timestamp
,`updated_at` datetime not null default current_timestamp on update current_timestamp
,`created_by` int not null
,`updated_by` int not null
,`organization_id` int not null
,`login_id` varchar(200) character set ascii
,`hashed_password` varchar(200) character set ascii
,`username` varchar(40)
,`role` varchar(20) character set ascii
,`provider` varchar(40) character set ascii
,`provider_id` varchar(40) character set ascii
,`provider_access_token` text character set ascii
,`provider_refresh_token` text character set ascii
,`removed` tinyint(1) not null
,primary key(`id`)
,unique(`organization_id`, `login_id`)
,foreign key(`organization_id`) references `organization`(`id`) on delete cascade
);
