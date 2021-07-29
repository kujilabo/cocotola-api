create table `space` (
 `id` integer primary key autoincrement
,`version` int not null
,`created_at` datetime not null default current_timestamp
,`updated_at` datetime not null default current_timestamp
,`created_by` int not null
,`updated_by` int not null
,`organization_id` int not null
,`type` int not null
,`key` varchar(20) not null
,`name` varchar(20) not null
,`description` varchar(40) not null
,unique(`organization_id`, `type`, `key`)
,foreign key(`created_by`) references `app_user`(`id`)
,foreign key(`updated_by`) references `app_user`(`id`)
,foreign key(`organization_id`) references `organization`(`id`)
);
