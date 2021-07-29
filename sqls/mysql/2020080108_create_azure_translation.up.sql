create table `azure_translation` (
 `text` varchar(30) character set ascii not null
,`lang` varchar(2) character set ascii not null
,`result` text not null
,primary key(`text`, `lang`)
);
