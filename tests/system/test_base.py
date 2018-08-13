from ciscoampbeat import BaseTest

import os


class Test(BaseTest):

    def test_base(self):
        """
        Basic test with exiting CiscoAMPBeats normally
        """
        self.render_config_template(
            path=os.path.abspath(self.working_dir) + "/log/*"
        )

        CiscoAMPBeats = self.start_beat()
        self.wait_until(lambda: self.log_contains("CiscoAMPBeats is running"))
        exit_code = CiscoAMPBeats.kill_and_wait()
        assert exit_code == 0
