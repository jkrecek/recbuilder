package template

const ENTITY_STRUCTURE_CODE = `
<?php
namespace {{NAMESPACE}};


use Krecek\Database\Annotation\Column;
use Krecek\Database\Annotation\OnCreate;
use Krecek\Database\Annotation\Primary;
use Nette\Utils\DateTime;

trait {{ENTITY_STRUCTURE_NAME}}
{
	{{ENTITY_VALUES}}
}
`
