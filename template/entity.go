package template

const ENTITY_CODE = `
<?php
namespace {{NAMESPACE}};

use Krecek\Database\Annotation\Column;
use Krecek\Database\Annotation\FormControl;
use Krecek\Database\Annotation\OnCreate;
use Krecek\Database\Annotation\OnUpdate;
use Krecek\Database\Annotation\Primary;
use Krecek\Database\Annotation\Table;
use Krecek\Database\StoredEntity;
use Nette\Utils\DateTime;

/**
 * {{READABLE_NAME}} entity.
 *
 * @Table("{{TABLE_NAME}}")
 */
class {{ENTITY_NAME}} extends StoredEntity
{
    use {{ENTITY_STRUCTURE_NAME}};

}
`
